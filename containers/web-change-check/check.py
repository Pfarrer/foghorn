import hashlib
import json
import os
import re
import sys
import time
from datetime import datetime, timezone
from urllib.parse import urlparse


def make_output(status, message, data, duration_ms):
    return {
        "status": status,
        "message": message,
        "data": data,
        "timestamp": datetime.now(timezone.utc).isoformat(),
        "duration_ms": round(duration_ms, 2),
    }


def output_and_exit(status, message, data, duration_ms):
    print(json.dumps(make_output(status, message, data, duration_ms)))
    sys.exit(0)


def parse_timeout(val):
    if not val:
        return 30
    s = val.rstrip("s")
    try:
        return int(s)
    except ValueError:
        return 30


def classify_error(error_message):
    msg = str(error_message).lower()
    if "ssl" in msg or "certificate" in msg or "cert" in msg:
        return "ssl"
    if "timeout" in msg or "timed out" in msg:
        return "timeout"
    if (
        "dns" in msg
        or "connect" in msg
        or "network" in msg
        or "refused" in msg
        or "unreachable" in msg
    ):
        return "network"
    return "browser"


def main():
    start = time.monotonic()

    url = os.environ.get("URL", "").strip()
    xpath_expr = os.environ.get("XPATH", "").strip() or None
    wait_seconds = int(os.environ.get("WAIT_SECONDS", "0").strip() or "0")
    wait_for_selector = os.environ.get("WAIT_FOR_SELECTOR", "").strip() or None
    verify_ssl = os.environ.get("VERIFY_SSL", "true").strip().lower() != "false"
    timeout = parse_timeout(os.environ.get("FOGHORN_TIMEOUT", "").strip())
    persistent_dir = os.environ.get("FOGHORN_PERSISTENT_DIR", "").strip()

    if not url:
        duration_ms = (time.monotonic() - start) * 1000
        output_and_exit(
            "fail", "URL environment variable is required", {"url": ""}, duration_ms
        )

    parsed = urlparse(url)
    if parsed.scheme not in ("http", "https"):
        duration_ms = (time.monotonic() - start) * 1000
        output_and_exit(
            "fail",
            f"Invalid URL scheme: {parsed.scheme}. Must be http or https.",
            {"url": url},
            duration_ms,
        )

    if not persistent_dir:
        duration_ms = (time.monotonic() - start) * 1000
        output_and_exit(
            "fail",
            "FOGHORN_PERSISTENT_DIR is required for state persistence",
            {"url": url},
            duration_ms,
        )

    if not os.path.isdir(persistent_dir):
        duration_ms = (time.monotonic() - start) * 1000
        output_and_exit(
            "fail",
            f"Persistent directory does not exist: {persistent_dir}",
            {"url": url},
            duration_ms,
        )

    state_file = os.path.join(persistent_dir, "state.json")

    try:
        from playwright.sync_api import sync_playwright
    except Exception as e:
        duration_ms = (time.monotonic() - start) * 1000
        output_and_exit(
            "unknown",
            f"Browser initialization failed: {e}",
            {"url": url, "error_type": "browser"},
            duration_ms,
        )

    content = None
    try:
        with sync_playwright() as p:
            browser = p.chromium.launch(
                headless=True,
                args=["--no-sandbox"] if not verify_ssl else [],
            )
            context = browser.new_context(
                ignore_https_errors=not verify_ssl,
            )
            page = context.new_page()
            page.set_default_timeout(timeout * 1000)

            try:
                page.goto(url, wait_until="domcontentloaded", timeout=timeout * 1000)
            except Exception as nav_err:
                browser.close()
                duration_ms = (time.monotonic() - start) * 1000
                error_type = classify_error(nav_err)
                output_and_exit(
                    "unknown",
                    str(nav_err),
                    {"url": url, "error_type": error_type},
                    duration_ms,
                )

            if wait_for_selector:
                try:
                    page.wait_for_selector(wait_for_selector, timeout=timeout * 1000)
                except Exception as wait_err:
                    browser.close()
                    duration_ms = (time.monotonic() - start) * 1000
                    error_type = classify_error(wait_err)
                    output_and_exit(
                        "unknown",
                        f"Wait for selector failed: {wait_err}",
                        {"url": url, "error_type": error_type},
                        duration_ms,
                    )

            if wait_seconds > 0:
                time.sleep(wait_seconds)

            if xpath_expr:
                try:
                    elements = page.locator(f"xpath={xpath_expr}").all()
                except Exception as xpath_err:
                    browser.close()
                    duration_ms = (time.monotonic() - start) * 1000
                    output_and_exit(
                        "fail",
                        f"Invalid XPath expression: {xpath_err}",
                        {"url": url, "xpath": xpath_expr},
                        duration_ms,
                    )

                if not elements:
                    browser.close()
                    duration_ms = (time.monotonic() - start) * 1000
                    output_and_exit(
                        "fail",
                        f"No elements matched XPath: {xpath_expr}",
                        {"url": url, "xpath": xpath_expr},
                        duration_ms,
                    )

                texts = []
                for el in elements:
                    t = el.text_content()
                    if t:
                        texts.append(t.strip())
                content = "\n".join(texts)
            else:
                content = page.locator("body").text_content()
                if content:
                    content = content.strip()

            browser.close()
    except Exception as e:
        duration_ms = (time.monotonic() - start) * 1000
        error_type = classify_error(e)
        output_and_exit(
            "unknown", str(e), {"url": url, "error_type": error_type}, duration_ms
        )

    if not content:
        duration_ms = (time.monotonic() - start) * 1000
        output_and_exit(
            "fail", "No content extracted from page", {"url": url}, duration_ms
        )

    current_hash = hashlib.sha256(content.encode("utf-8")).hexdigest()

    data = {
        "url": url,
        "changed": False,
        "current_hash": current_hash,
    }
    if xpath_expr:
        data["xpath"] = xpath_expr

    if not os.path.exists(state_file):
        state = {
            "hash": current_hash,
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "url": url,
        }
        with open(state_file, "w") as f:
            json.dump(state, f, indent=2)
        data["first_run"] = True
        duration_ms = (time.monotonic() - start) * 1000
        output_and_exit("pass", "Initial snapshot saved", data, duration_ms)

    try:
        with open(state_file, "r") as f:
            previous_state = json.load(f)
    except Exception:
        state = {
            "hash": current_hash,
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "url": url,
        }
        with open(state_file, "w") as f:
            json.dump(state, f, indent=2)
        data["first_run"] = True
        duration_ms = (time.monotonic() - start) * 1000
        output_and_exit("pass", "Initial snapshot saved", data, duration_ms)

    previous_hash = previous_state.get("hash", "")

    if current_hash == previous_hash:
        state = {
            "hash": current_hash,
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "url": url,
        }
        with open(state_file, "w") as f:
            json.dump(state, f, indent=2)
        data["changed"] = False
        duration_ms = (time.monotonic() - start) * 1000
        output_and_exit("pass", "No change detected", data, duration_ms)
    else:
        state = {
            "hash": current_hash,
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "url": url,
        }
        with open(state_file, "w") as f:
            json.dump(state, f, indent=2)
        data["changed"] = True
        data["previous_hash"] = previous_hash
        duration_ms = (time.monotonic() - start) * 1000
        output_and_exit("fail", "Change detected on page", data, duration_ms)


if __name__ == "__main__":
    main()
