# ansible-unused-role-finder

`ansible-unused-role-finder` lists Ansible roles that no playbook uses, so you can mechanically find roles left behind after a system is replaced.

## How roles are resolved

A role counts as used when a playbook declares it at the play level, or when it is pulled in as a `meta/main.yml` dependency of a used role. This matches the role layout in [ansible-role-practice](https://github.com/zinrai/ansible-role-practice), where a role maps to one package and is the minimum unit of setup.

## Limitations

Roles reached through `include_role` or `import_role` are not analyzed, because they are resolved at task runtime rather than declared statically. A role used only that way can therefore be reported as unused. Enforcing that roles stay declared at the play level is out of scope here.

## Usage

Pass `-playbook` once per playbook that shares the roles directory. A role is unused only when none of the given playbooks reach it, directly or through dependencies, so playbooks that share roles must be checked together.

```bash
./ansible-unused-role-finder -playbook consul-server.yml -playbook nomad-server.yml -roles roles
```

### Flags

- `-playbook`: Path to an Ansible playbook YAML file, repeatable (required)
- `-roles`: Path to the Ansible roles directory (required)

### Exit status

- `0`: the check ran, and any unused roles were printed to stdout
- `1`: an error occurred

## Output

Unused role names go to stdout, one per line, sorted. No output means every role is used.

## License

This project is licensed under the [MIT License](./LICENSE).
