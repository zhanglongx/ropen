# ropen

ropen is a simple http server that provides quick access to directories and files.

ropen access to directories and files and outputs this url to stdout. this feature is particularly simple and convenient in ssh environments, where a url can be detected by modern terminal software (e.g., Windows Terminal) the user can click on the link, open a browser, so quickly access the remote directories and files.

Modern browsers are usually very careful when downloading files and will usually block file download operations if the connection is not https or if the https certificate is not available. Therefore ropen provides a ca function that generates a ca-signed integer.

⚠️If user use self-signed ca certificate, they need to install the ca certificate in the operating system.

## Usage

```shell
ropen [options] [path]
```

## Configuration

### Searching Order

ropen searches for configuration files in the following order:

- user specific configuration file via `-cfg`

- $(pwd)/ropen.yml

- $HOME/.ropen.yml

### Configuration File

```yml
port: [38080]

preferips: 
  - <ip>

issuer:
  capath: <path to issuer crt>
  keypath: <path to issuer key>
```

- preferips: specify the ip address to bind to.

- issuer: ca in crt and pkcs8 format.

- fields can be omitted. ropen will use the default value.

## TODO

- [ ] support ipv6
