# Mr. Lovenstein RSS

A small app to rebuild <https://tapas.io/rss/series/3346> into an RSS feed compatible with IFTTT.

Feed is forwarded to <https://t.me/iftttmrlovenstein>.

## Example

Input :

```xml
<?xml version="1.0" encoding="UTF-8"?>
<rss xmlns:content="http://purl.org/rss/1.0/modules/content/" version="2.0">
  <channel>
    <title>Mr. Lovenstein</title>
    <link>https://tapas.io/series/3346</link>
    <description>World's Sweatiest Comic</description>
    <lastBuildDate>Fri, 06 Aug 2021 12:31:05 GMT</lastBuildDate>
    <image>
      <title>Mr. Lovenstein</title>
      <url>https://d30womf5coomej.cloudfront.net/sa/b6/61a71c65-af8c-4997-8a6a-9987d6678c5d.png</url>
      <width>-1</width>
      <height>-1</height>
    </image>
    <item>
      <title>Hi Strung</title>
      <link>https://tapas.io/episode/2253285</link>
      <content:encoded>&lt;p&gt;I will never recover from this.&lt;/p&gt;&lt;img src="https://d30womf5coomej.cloudfront.net/sa/90/fa96aa67-5ab4-4eca-9424-32ded1512f78.png"/&gt;</content:encoded>
      <pubDate>Fri, 06 Aug 2021 12:31:05 GMT</pubDate>
      <author>J. L. Westover</author>
    </item>
    <item>
      <title>Who's Gonna Call</title>
      <link>https://tapas.io/episode/2250364</link>
      <content:encoded>&lt;p&gt;There's more than one way to bust a ghost.&lt;/p&gt;&lt;img src="https://d30womf5coomej.cloudfront.net/sa/37/14da2d15-065c-4599-b96d-0769c3f50233.png"/&gt;</content:encoded>
      <pubDate>Tue, 03 Aug 2021 13:12:09 GMT</pubDate>
      <author>J. L. Westover</author>
    </item>
  </channel>
</rss>
```

Output :

```xml
<?xml version="1.0" encoding="UTF-8"?><rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
  <channel>
    <title>Mr. Lovenstein</title>
    <link>https://tapas.io/rss/series/3346</link>
    <description>World's Sweatiest Comic</description>
    <pubDate>Fri, 06 Aug 2021 12:31:05 +0000</pubDate>
    <lastBuildDate>Fri, 06 Aug 2021 12:31:05 +0000</lastBuildDate>
    <item>
      <title>Hi Strung.</title>
      <link>https://tapas.io/episode/2253285</link>
      <description></description>
      <content:encoded><![CDATA[https://d30womf5coomej.cloudfront.net/c/4d/1b90af4e-98e7-46fb-ab3c-a67a52d88582.png]]></content:encoded>
      <author>J. L. Westover</author>
      <pubDate>Fri, 06 Aug 2021 12:31:05 +0000</pubDate>
    </item>
    <item>
      <title>Who's Gonna Call.</title>
      <link>https://tapas.io/episode/2250364</link>
      <description></description>
      <content:encoded><![CDATA[https://d30womf5coomej.cloudfront.net/c/23/17a58e79-a8f2-448f-b690-4a94e40931a0.png]]></content:encoded>
      <author>J. L. Westover</author>
      <pubDate>Tue, 03 Aug 2021 13:12:09 +0000</pubDate>
    </item>
  </channel>
</rss>
```

## License

Unless explicitly stated to the contrary, all contents licensed under the [MIT License](LICENSE).

Comics are the sole property of J. L. Westover. More info on <https://tapas.io/series/MrLovenstein/info>.
