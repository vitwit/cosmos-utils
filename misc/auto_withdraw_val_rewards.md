# Auto withdraw validator rewards

The script withdraws the validator rewards everymonth for given array of validators using a config. It leverages authz features for doing these transactions.

### In this script:

The config.json file is read using the jq tool, which is a command-line JSON processor.

The script then iterates through the networks specified in the JSON configuration, extracting the parameters for each network and performing the desired operations.

Make sure to replace the ellipses (...) in the JSON configuration with your specific network parameters. You can add more network configurations to the `auto_rewards_config.json`
file as needed, and the script will process them all in a loop.

### Crontab for running every month at specified time.
To run the script every month on the 1st day at 9 PM indefinitely, you can use the cron scheduler on a Unix-based system like Linux or macOS. Here's how you can set up a cron job to achieve this:

Open your terminal.

Edit your user's crontab file using the following command:

```bash
crontab -e
```
Add the following line to schedule your script:
```bash
0 21 1 * * /vitwit/cosmos-utils/scripts/auto_rewards.sh
```

Where
- 0 represents the minute (0-59).
- 21 represents the hour (24-hour format).
- 1 represents the day of the month (1-31).
- `*` for the month (1-12).
- `*` for the day of the week (0-6, where 0 is Sunday).
- `/vitwit/cosmos-utils/scripts/auto_rewards_script.sh` should be replaced with the actual path to your shell script path.
In this example, the script will run every month on the 1st day at 9 PM (21:00).

Save and exit the editor.

The cron job is now set up to run your script at the specified time and frequency. Be sure to replace /path/to/your/script.sh with the actual path to your shell script.




