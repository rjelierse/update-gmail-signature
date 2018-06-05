# Update GMail signature

> A simple utility to update the GMail signature of all users in your G Suite organization

## Installation

```
$ go get github.com/rjelierse/update-gmail-signature 
```

## Usage

To use this program, a valid Google API service account is required. This account should have
[delegated domain-wide authority](https://developers.google.com/admin-sdk/directory/v1/guides/delegation)
to allow changing user details.

You should allow access to these scopes:

* `https://www.googleapis.com/auth/admin.directory.user.readonly`
* `https://www.googleapis.com/auth/gmail.settings.basic`

Also, you'll need an HTML template to use as a basis for the email signature.

### Flags

* `-secret`: Set the path the the API credentials file. Defaults to `client_secret.json` in the directory
  where the command is executed.
* `-template`: Set the path to the signature template. Defaults to `template.html` file in the directory
  where the command is executed.
* `-domain`: The organization to use when looking up users in the G Suite directory.
* `-subject`: The user to impersonate when looking up users in the G Suite directory.
  This user should have full access to the directory.

### Template variables

The following variables are available (for now):

* `Name`: the full name of the user
* `Title`: the job title in the primary organization
* `Mobile`: the phone number marked as 'mobile'
* `Address`: the address marked as 'work'
* `Phone`: the phone number marked as 'work'
