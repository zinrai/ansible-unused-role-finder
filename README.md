# ansible-unused-role-finder

`ansible-unused-role-finder` is a command-line tool written in Go that helps identify unused Ansible roles in your playbooks. It analyzes your Ansible playbook and roles directory to determine which roles are not being used, either directly or as dependencies of other roles.

I wanted to mechanically list Ansible Roles that were no longer needed when the system was replaced.

## Features

- Identifies unused Ansible roles
- Handles role dependencies
- Simple output format, listing only unused roles

## Installation

```bash
$ go build
```

## Usage

Run the tool using the following command:

```bash
./ansible-unused-role-finder -playbook /path/to/your/playbook.yml -roles /path/to/your/roles/directory
```

### Flags

- `-playbook`: Path to the Ansible playbook YAML file (required)
- `-roles`: Path to the Ansible roles directory (required)

### Example

```bash
./ansible-unused-role-finder -playbook /etc/ansible/site.yml -roles /etc/ansible/roles
```

## Output

The tool will output a list of unused roles, one per line. If there are no unused roles, it will not produce any output.

Example output when unused roles are found:

```
unused_role1
unused_role2
unused_role3
```

If all roles are used, there will be no output.

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
