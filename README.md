alertmanager-webhook-free is a gateway that receives alertmanager webhooks and send the alert to Free mobile API

# Requirements

- Go (version >= 1.15)
- Alertmanager (version >= 0.21)

# Build
`$ go build`

# Usage
Enable SMS Notification through SMS on your account : https://mobile.free.fr/moncompte

Use the -h flag to see full usage:

```
$ alertmanager-webhook-free -h
Usage of alertmanager-webhook-free:
  -config string
        Path to config
```

Config is writed in yaml
```
sentry:
  dsn: "https://sentry.io"
server:
  address: ":9234"
free:
  user: "xxxxxxx"
  pass: "yyyyyyy"
```

`$ alertmanager-webhook-free -config /etc/alertmanager-webhook-free/alertmanager-webhook-free.yml`

# Licence

The code is under CeCILL license.

You can find all details here: https://cecill.info/licences/Licence_CeCILL_V2.1-en.html

# Credits

Copyright © Ludovic Ortega, 2021

Contributor(s):

-Ortega Ludovic - ludovic.ortega@adminafk.fr