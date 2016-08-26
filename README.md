# LogrusFluentdTaggedUdp
[Logrus Hook](https://github.com/Sirupsen/logrus) for the Fluentd [Tagged UDP plugin](https://github.com/toyokazu/fluent-plugin-tagged_udp)

This sends to Fluentd over UDP in the format of [tag][delimiter][JSON data] where the default delimiter is a tab.  This is the format the the Tagged UDP plugin expects.
