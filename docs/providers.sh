#/usr/bin/env bash

# This script generates a markdown file containing a list of all providers and
# supported objects.

set -e
PROJECT_ROOT=$(cd "$(dirname $0)/.." && pwd)
cd $PROJECT_ROOT

FRIEZA="$PROJECT_ROOT/cmd/frieza/frieza"
if ! [ -f "$FRIEZA" ]; then
    echo "$FRIEZA not found, abort"
    exit 1
fi

output_file=$PROJECT_ROOT/docs/providers.md
echo "# Providers and supported objects" > $output_file
$FRIEZA provider list | while read provider; do
    echo "" >> $output_file
    echo "## $provider" >> $output_file
    $FRIEZA provider describe $provider | while read object_name; do
        echo "- $object_name" >> $output_file
    done
done