source "$PWD"/scripts/common.sh

# Init the provider
/bin/bash "$PWD"/scripts/ccv-provider.sh
# Init the consumer
/bin/bash "$PWD"/scripts/ccv-consumer.sh

echo
cecho "YELLOW" "To init sample:"
cecho "GREEN" "./scripts/init-sample.sh"
echo
cecho "YELLOW" "To track hermes:"
cecho "GREEN" "tmux a -t hermes"
echo
cecho "YELLOW" "To track mantrachain:"
cecho "GREEN" "tmux a -t mantrachain"
echo
cecho "YELLOW" "To track provider:"
cecho "GREEN" "tmux a -t provider"
echo
