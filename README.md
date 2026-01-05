# creds-forwarder
Forwarding credentials like ssh-agent

Handy tool to be able to forward cloud credentials like how one forwards their SSH credentials over a SSH connection.

Combined with the plain old SSH agent forwarding, only your laptop stays logged into your cloud accounts, and all your remote machines are secret-less.

## Usage (AWS)

1. Make sure you are logged into AWS.
2. Start the forwarding utility `aws configure export-credentials --profile qsg-dogfooding-nonprod.write | go run ./`
3. Assume it serves at `/tmp/auth.sock`, now you can SSH to your remote machine with an additional port forwarding config `ssh user@remote -R /tmp/auth.sock:/tmp/auth.sock`.
4. Make sure your remote AWS config `~/.aws/config` looks like the following:
```
[default]
credential_process = curl --unix-socket /tmp/awsauth.sock http://./>
```
5. Profit! Actually, SSH doesn't clean up the forwarded used socket. So adding this hook to the end of your shell rc (.bashrc/.zshrc/etc) script helps: `function onexit { rm -f /tmp/auth.sock; }; trap onexit EXIT`.

## Usage (GCP)

Coming soon...
