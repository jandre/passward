# Introducing passward

Passward is a open-source password sharing app designed to be used by teams.

# Why use passward?

1. Transparency: It doesn't require your 3rd passwords to be hosted by a third party behind some hidden/dubious security controls.  All the source code is available here.
2. Auditability: It allows you to control and log changes.
3. Distributed: It is distributed using `git` which we all know and love.
4. Sharing: You can easily share passwords with other members of your team.  It's as simple as giving them 
access to the git repo and running `passward share <vault name>`

# Design

1. Passward configs are stored in ~/.passward/

2. Passward vaults are stored in ~/.passward/vault/<name>.  Each vault corresponds to a git repo.

3. Vaults are organized as follows:

~/.passward/vault/<name>

```
users/
  bob/
    keys/
      home_id_pub.rsa
    encrypted_master


config/
  index
  encrypted_master
keys/
  com.blah.test/
    user
    passphrase
    description
  myssh/
    user
    passphrase
    description
    ...
```

# Q&A

*Q. How do I add read-only users?*

A delightful question! Simply give the user read-only access to the 
remote repository to the vault.

*Q. Should I store my passwords in Github, even if they are encrypted?*

Probably not.  You should use a private git server if you can.  







