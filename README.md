# Homebase

Homebase grabs your current public IP address and updates a Digital Ocean hostname
to point to the IP.  You can use this to provide your own dynamic DNS for home
computers.

## Installation

Download the [latest release](https://github.com/bsdavidson/homebase/releases/latest) and rename the file to `homebase`.

## Usage

Run from commandline:

    $ homebase --domain example.com --record subdomain --token YOUR_API_TOKEN

Your domain will need to be hosted with Digital Ocean and you'll need a
[Digital Ocean API token](https://cloud.digitalocean.com/settings/api/tokens).

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/bsdavidson/homebase.

## License

The project is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).

