# Update GMail signature

> A simple utility to update the GMail signature of all users in your G Suite organization

## Installation

```
$ go get github.com/rjelierse/update-gmail-signature 
```

## Usage

To use this program, valid Google API credentials are required. You'll need OAuth client credentials,
downloaded in JSON format.

Also, you'll need an HTML template to use as a basis for the email signature.

### Flags

* `-secret`: Set the path the the API credentials file. Defaults to `client_secret.json` in the directory
  where the command is executed.
* `-template`: Set the path to the signature template. Defaults to `template.html` file in the directory
  where the command is executed.

### Template variables

The following variables are available (for now):

* `Name`: the full name of the user
* `Title`: the job title in the primary organization
* `Mobile`: the phone number marked as 'mobile'
