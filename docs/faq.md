# FAQ

**Q: How do I update from version 5 to version 6**

A: First, don't use the update.js script, it won't work since the server portion of version 6 as been completely rewritten in GO. Simply unzip the archive that contains the FreeScanner executable and its PDF document to a new folder, then copy the database.sqlite from version 5 to the new folder that contains version 6 and make sure to rename it to _freescanner.db_.

**Q: I tried the autocert function but I get strange error messages**

A: Due to the ACME protocol used by Let's Encrypt, ports 80 and 443 must be open to the world for the autocert to work. The domain specified via the `-ssl_auto_cert` argument must also match the IP address of your FreeScanner instance.

**Q: The web app keeps displaying a dialog telling me that a new version is available**

A: Force a refresh of the web application from the browser (usually with ctrl-shift-r) to resolve the issue. Alternatively, you can click on the icon just to the left of the URL address and select website settings, then clear all website data.

**Q: How do I install FFMPEG on Windows**

A: Please follow instructions at this address: [https://www.wikihow.com/Install-FFmpeg-on-Windows](https://www.wikihow.com/Install-FFmpeg-on-Windows)

**Q: How do I configure a reverse-proxy in front of FreeScanner**

A: There are so many reverse proxy technologies out there that it's hard the cover them all. One thing to keep in mind is that FreeScanner works with websockets, so the reverse proxy must also supports websockets to work properly with FreeScanner. For some examples, take a look at the [https://github.com/amigan/freescanner/tree/master/docs/examples/apache](https://github.com/amigan/freescanner/tree/master/docs/examples/apache) for `Apache HTTP` or [https://github.com/amigan/freescanner/tree/master/docs/examples/nginx](https://github.com/amigan/freescanner/tree/master/docs/examples/nginx) for `nginx`.

**Q: How do I get notified when a new release is available**

A: Use the GitHub `watch` feature. This requires you to have a GitHub account, which you can create for free. Go to the FreeScanner repository at [https://github.com/amigan/freescanner](https://github.com/amigan/freescanner) and select the `watch` button. You can be notified of every change made to the repository, or simply be notified when a new release is available.

**Q: How can I listen to multiple instances from the same server**

A: Simply open a new browser tab to the same URL with a special `id` parameter that will distinguish each instance from the other. This allows you to remember the selection of talkgroups for each of the instances. Without the `id` parameter, only the last talkgroups selection is remembered across all instances. For example: `http://localhost:3000/?id=instance2`.

**Q: I did not find an answer to my question in this FAQ**

A: No problem, just drop us a line at [freescanner@saubeo.solutions](mailto:freescanner@saubeo.solutions) and we'll make sure to add the relevant information in this document in the next release. In the meantime, You can ask your questions on the [FreeScanner Discussions](https://github.com/amigan/freescanner/discussions) at [https://github.com/amigan/freescanner/discussions](https://github.com/amigan/freescanner/discussions).

\pagebreak{}