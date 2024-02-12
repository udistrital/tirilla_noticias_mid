
#!/usr/bin/env bash

set -e
set -u
set -o pipefail

#if [ -n "${PARAMETER_STORE:-}" ]; then
#  export NOTICIAS_MID_PGUSER="$(aws ssm get-parameter --name /${PARAMETER_STORE}noticias_mid/db/username --output text --query Parameter.Value)"
#  export NOTICIAS_MID_PGPASS="$(aws ssm get-parameter --with-decryption --name /${PARAMETER_STORE}/noticias_mid/db/password --output text --query Parameter.Value)"

exec ./main "$@"