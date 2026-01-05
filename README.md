# creds-forwarder
Forwarding credentials like ssh-agent

Handy tool to be able to forward cloud credentials like how one forwards their SSH credentials over a SSH connection.

Combined with the plain old SSH agent forwarding, only your laptop stays logged into your cloud accounts, and all your remote machines are secret-less.

## Usage (AWS)

1. Make sure you are logged into AWS. Install the utility via `go install github.com/yuedongze/creds-forwarder@latest`.
2. Start the forwarding utility `aws configure export-credentials --profile <profile-name> | creds-forwarder`.
3. Assume it serves at `/tmp/auth.sock`, now you can SSH to your remote machine with an additional port forwarding config `ssh user@remote -R /tmp/auth.sock:/tmp/auth.sock`.
4. Make sure your remote AWS config `~/.aws/config` looks like the following:
```
[default]
credential_process = curl --silent --unix-socket /tmp/auth.sock http://./token
```
5. Profit! Try running AWS commands on remote like `aws sts get-caller-identity`.
6. Actually, SSH doesn't clean up the forwarded used socket. So adding this hook to the end of your shell rc (.bashrc/.zshrc/etc) script helps: `function onexit { rm -f /tmp/auth.sock; }; trap onexit EXIT`.

## Usage (GCP)

Coming soon...
