# ipinfo

> A poor-mans IPAM alternative

1. Add all of your AWS accounts to accounts.yml and copy it to `$HOME/.accounts.yml`
2. Get information about an ip address with `ipinfo 192.168.100.2`

ipinfo will then perform a few things:

1. Lookup the account name from the .accounts.yml file
2. Use the account name as AWS_PROFILE when authenticating against AWS
3. Find a Network Interface with the given IP
4. Print a description if set
5. Get a Name tag of the associated EC2 instance if an instanceID is found
