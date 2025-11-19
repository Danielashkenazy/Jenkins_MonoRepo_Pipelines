CHANGED_FILES=$(git diff --name-only HEAD~1 HEAD)

SERVICES_CHANGED=""

for SERVICE in user-service transaction-service notification-service; do
    if echo "$CHANGED_FILES" | grep -q "^$SERVICE/"; then
        SERVICES_CHANGED="$SERVICES_CHANGED $SERVICE"
    fi
done

echo $SERVICES_CHANGED
