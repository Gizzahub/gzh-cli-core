#!/bin/sh
# Exercise the real Makefile with a controlled formatter, without editing Go files.
set -eu
mkdir -p tmp
fixture=$(mktemp -d tmp/fmt-check.XXXXXX)
cleanup() {
    rm -f "$fixture/gofmt" "$fixture/stdout" "$fixture/stderr"
    rmdir "$fixture"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
cat > "$fixture/gofmt" <<'FORMATTER'
#!/bin/sh
test "$#" -eq 2 && test "$1" = -l && test "$2" = . || exit 99
case "$FMT_CHECK_TEST_MODE" in
    clean) exit 0 ;;
    unformatted) echo example.go ;;
    error) echo 'injected gofmt diagnostic' >&2; exit 23 ;;
    error-with-output) echo example.go; echo 'injected gofmt diagnostic' >&2; exit 23 ;;
    *) exit 99 ;;
esac
FORMATTER
chmod +x "$fixture/gofmt"
for mode in clean unformatted error error-with-output; do
    status=0
    PATH="$(pwd)/$fixture:$PATH" FMT_CHECK_TEST_MODE="$mode" \
        make --no-print-directory fmt-check > "$fixture/stdout" 2> "$fixture/stderr" || status=$?
    if [ "$mode" = clean ]; then
        test "$status" -eq 0
        grep -q 'Code is properly formatted' "$fixture/stdout"
    else
        test "$status" -ne 0
        if grep -q 'Code is properly formatted' "$fixture/stdout"; then
            echo "fmt-check falsely reported success for $mode" >&2
            exit 1
        fi
        case "$mode" in
            unformatted) grep -q example.go "$fixture/stdout" ;;
            error*) grep -q 'injected gofmt diagnostic' "$fixture/stderr" ;;
        esac
    fi
    printf 'PASS fmt-check %s\n' "$mode"
done
