import json
import os
import matplotlib.pyplot as plt
from datetime import datetime

def main():
    # Define the path to the genesis file
    genesis_file = os.path.join('dukong', 'genesis.json')

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
        amount = int(balance['coins'][0]['amount'])
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

    # Recursively extract all nested parameters
    def extract_params(data, depth=0):
        result = ""
        indent = "  " * depth
        if isinstance(data, dict):
            for key, value in data.items():
                if isinstance(value, (dict, list)):
                    result += f"{indent}- **{key.replace('_', ' ').title()}**:\n"
                    result += extract_params(value, depth + 1)
                else:
                    result += f"{indent}- **{key.replace('_', ' ').title()}**: {value}\n"
        elif isinstance(data, list):
            for index, item in enumerate(data):
                result += f"{indent}- **Item {index + 1}**:\n"
                result += extract_params(item, depth + 1)
        return result

    # Generate Markdown report with enhanced styling and all parameters
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

        f.write('## Full Application State and Parameters\n')
        f.write(extract_params(genesis['app_state']))
        f.write('\n')

        f.write('## Full Consensus Parameters\n')
        f.write(extract_params(genesis['consensus']['params']))
        f.write('\n')

        f.write('---\n')
        f.write(f'Report generated on {datetime.utcnow().strftime("%Y-%m-%d %H:%M:%S UTC")}\n')

    print("Report generated: report.md")

if __name__ == '__main__':
    main()
