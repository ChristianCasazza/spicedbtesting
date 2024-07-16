#!/bin/bash

FILE="relationships.zed"
PERMISSIONS_SYSTEM="oceanprotocol_testing"

while IFS= read -r line
do
  zed relationship create $line --permissions-system $PERMISSIONS_SYSTEM
done < "$FILE"
