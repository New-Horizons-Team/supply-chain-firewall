# iFood Supply-Chain Firewall

iFood Supply-Chain Firewall is a command-line tool for preventing the installation of malicious PyPI and go packages.  It is intended primarily for use by engineers to protect their development workstations from compromise in a supply-chain attack.

Given a command for a supported package manager, iFood Supply-Chain Firewall collects all package targets that would be installed by the command and checks them against reputable sources of data on open source malware and vulnerabilities.  The command is automatically blocked from running when any data source finds that any target is malicious.  In cases where a data source reports other findings for a target, they are presented to the user along with a prompt confirming intent to proceed with the installation.

Default data sources include:

- Datadog Security Research's public [malicious packages dataset](https://github.com/DataDog/malicious-software-packages-dataset)
- [OSV.dev](https://osv.dev) advisories

The principal goal of iFood Supply-Chain Firewall is to block 100% of installations of known-malicious packages within the purview of its data sources.

## Getting started

### Installation

Go to the Releases (TODO) page, download the most recent version for your system, and install it to a directory available in your path (e.g., `~/.local/bin`).

### Post-installation steps

To get the most out of iFood Supply-Chain Firewall, it is recommended to run the `scfw configure` command after installation.  This script will walk you through configuring your environment so that all commands for supported package managers are passively run through `scfw` as well as enabling Datadog logging, described in more detail below.

```bash
$ scfw configure
...
```

### Compatibility and limitations

|  Package manager  |  Compatible versions  |        Inspected subcommands                        |
| :---------------: | :-------------------: | :-------------------------------------------------: |
| pip               | >= 22.2               | `install`                                           |
| poetry            | >= 1.7                | `add`, `install`, `sync`, `update`                  |
| go                | >= 1.17.0             | `build`, `generate`, `get`, `install`, `mod`, `run` |

In keeping with its goal of blocking 100% of known-malicious package installations, `scfw` will refuse to run with an incompatible version of a supported package manager.  Please upgrade to or verify that you are running a compatible version before using this tool.

iFood Supply-Chain Firewall may only know how to inspect some of the "installish" subcommands for its supported package managers.  These are shown in the above table.  Any other subcommands are always allowed to run.

Currently, iFood Supply-Chain Firewall is fully supported on Linux and macOS systems, though it should also run as intended on Windows.

### Uninstalling iFood Supply-Chain Firewall

iFood Supply-Chain Firewall may by simply removing it from where it was installed.  Before doing so, be sure to run the command `scfw configure --remove` to remove any iFood Supply-Chain Firewall-managed configuration you may have previously added to your environment.

```bash
$ scfw configure --remove
...
```

## Usage

To use iFood Supply-Chain Firewall, simply prepend `scfw run --` to the command you want to run.

```
$ scfw run -- pip install -r requirements.txt
$ scfw run -- poetry add git+https://github.com/DataDog/guarddog
$ scfw run -- go mod download
```

For `pip install` commands, packages will be installed in the same environment (virtual or global) in which the command was run.

## Acknowledgments

iFood Supply-Chain Firewall was initially implemented as a fork of [Supply-Chain Firewall](https://github.com/DataDog/supply-chain-firewall), but it's since been ported to Go.
