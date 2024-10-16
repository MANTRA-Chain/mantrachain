import json
import os
import matplotlib.pyplot as plt
from datetime import datetime
import matplotlib.ticker as mtick

def main():
    # Define the path to the genesis file
    genesis_file = os.path.join('mantra-1', 'genesis.json')

    # Load genesis.json
    with open(genesis_file, 'r') as f:
        genesis = json.load(f)

    # Extract general information
    general_info = {
        'App Name': genesis.get('app_name'),
        'App Version': genesis.get('app_version'),
        'Genesis Time': genesis.get('genesis_time'),
        'Chain ID': genesis.get('chain_id'),
        'Initial Height': genesis.get('initial_height'),
    }

    # Extract token information
    bank = genesis['app_state']['bank']
    total_supply = bank.get('supply', [])
    denom_metadata = bank.get('denom_metadata', [])
    balances = bank.get('balances', [])

    # Extract balances for plotting
    balance_data = []
    for balance in balances:
        address = balance['address']
        coins = balance.get('coins', [])
        if coins:  # Check if the coins list is not empty
            amount = int(coins[0]['amount'])
            balance_data.append({'address': address, 'amount': amount})

    # Sort balance data by amount
    balance_data.sort(key=lambda x: x['amount'], reverse=True)

    # Create a pie chart of the token distribution
    labels = [item['address'] for item in balance_data[:5]] + ['Others']
    sizes = [item['amount'] for item in balance_data[:5]]
    others = sum(item['amount'] for item in balance_data[5:])
    sizes.append(others)

    fig1, ax1 = plt.subplots()
    ax1.pie(sizes, labels=labels, autopct='%1.1f%%', startangle=140, colors=plt.cm.Paired.colors)
    ax1.axis('equal')  # Equal aspect ratio ensures that pie is drawn as a circle.
    plt.title('Token Distribution')
    plt.savefig('token_distribution.png')
    plt.close()

    # Extract validators
    gen_txs = genesis['app_state']['genutil']['gen_txs']
    validators = []
    for tx in gen_txs:
        for msg in tx['body']['messages']:
            if msg['@type'] == '/cosmos.staking.v1beta1.MsgCreateValidator':
                validators.append({
                    'Moniker': msg['description']['moniker'],
                    'Validator Address': msg['validator_address'],
                    'Self-Delegation': msg['value']['amount'] + ' ' + msg['value']['denom'],
                })

    # Extract governance parameters
    gov_params = genesis['app_state']['gov']['params']

    # Extract market map
    market_map = genesis['app_state']['marketmap']['market_map']['markets']

    # Generate Markdown report with enhanced styling
    with open('report.md', 'w') as f:
        f.write('# Genesis Report\n\n')
        f.write('![Company Logo](mantra.png){ width=150px }\n\n')

        f.write('## General Information\n')
        for key, value in general_info.items():
            f.write(f'**{key}**: {value}\n\n')

        f.write('## Token Information\n')
        f.write('### Total Supply\n')
        for supply in total_supply:
            f.write(f'- **{supply["denom"]}**: {int(supply["amount"]):,}\n')
        f.write('\n')

        f.write('### Denominations\n')
        for denom in denom_metadata:
            f.write(f'- **{denom["display"]}** ({denom["symbol"]}): {denom["description"]}\n')
        f.write('\n')

        f.write('## Token Distribution\n')
        f.write('![Token Distribution](token_distribution.png)\n\n')

        f.write('## Validators\n')
        for val in validators:
            f.write(f'### {val["Moniker"]}\n')
            f.write(f'- Validator Address: `{val["Validator Address"]}`\n')
            f.write(f'- Self-Delegation: {int(val["Self-Delegation"].split()[0]):,} {val["Self-Delegation"].split()[1]}\n\n')

        f.write('## Governance Parameters\n')
        for key, value in gov_params.items():
            # Format durations
            if 'period' in key or 'duration' in key:
                value = format_duration(value)
            # Format percentages
            elif 'threshold' in key or 'quorum' in key or 'ratio' in key:
                value = f"{float(value)*100:.2f}%"
            f.write(f'- **{key.replace("_", " ").title()}**: {value}\n')
        f.write('\n')

        f.write('## Market Map\n')
        for market, data in market_map.items():
            f.write(f'### {market}\n')
            f.write(f'- Decimals: {data["ticker"]["decimals"]}\n')
            f.write(f'- Enabled: {"Yes" if data["ticker"]["enabled"] else "No"}\n')
            f.write('- Providers:\n')
            for provider in data.get('provider_configs', []):
                f.write(f'  - **{provider["name"]}**: {provider["off_chain_ticker"]}\n')
            f.write('\n')

        f.write('---\n')
        f.write(f'Report generated on {datetime.utcnow().strftime("%Y-%m-%d %H:%M:%S UTC")}\n')

    print("Report generated: report.md")

def format_duration(duration_str):
    # Assuming duration_str is in the format '3600s'
    if duration_str.endswith('s'):
        total_seconds = int(duration_str[:-1])
        hours = total_seconds // 3600
        minutes = (total_seconds % 3600) // 60
        seconds = total_seconds % 60
        return f"{hours}h {minutes}m {seconds}s"
    return duration_str

if __name__ == '__main__':
    main()

