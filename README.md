# Forex Board
The goal of the project is to provide monitoring of forex accounts.

## Using
1. Download the latest version at release page.
1. Unpack the archive to your directory.
1. Create separate file for each monitored account, for example `my-account.json`:
```json
{
  "name": "My Account",
  "location": "https://www.myfxbook.com/members/user1/account1/id",
  "refresh_seconds": 300,
  "currency_divider": 100
}
```
1. Run file:
```shell script
./fxbd
```
