// ----------------------------------------------------------------------
// 获取Notion页面的HTML
// 2020年5月6日10:20:21
// ----------------------------------------------------------------------
// GITHUB: https://github.com/yiningv
// ----------------------------------------------------------------------
'use strict';

const puppeteer = require("puppeteer");

async function getPageHtml(url) {
    const browser = await puppeteer.launch();
    const page = await browser.newPage();
    await page.goto(url);
    await page.waitForSelector('#notion-app');
    await page.waitFor(8000);
    const data = await page.evaluate(() => {
        const document = window.document;
        // 图片链接转换
        document.querySelectorAll('div.notion-page-content  img').forEach(item => {
            if (item.src.startsWith("https://s3.us-west")) {
                let [parsedOriginUrl] = item.src.split("?");
                item.src = `https://notion.so/image/${encodeURIComponent(parsedOriginUrl).replace("s3.us-west", "s3-us-west")}`;
            }
        })

        // TOC 链接转化
        const qs = "#notion-app > div > div.notion-cursor-listener > div > div.notion-scroller.vertical.horizontal > div.notion-page-content > div > div:nth-child(1) > div > a"
        document.querySelectorAll(qs).forEach(item => {

            const toDashID = (id) => {
                if (!id.match("^[a-zA-Z0-9]+$")) {
                    return id;
                }
                const res = `${id.substr(0, 8)}-${id.substr(8, 4)}-${id.substr(12, 4)}-${id.substr(16, 4)}-${id.substr(20, 32)}`;
                return res;
            }

            let u
            try {
                u = new URL(item.href)
            } catch (error) {
                console.log(error)
            }

            if (u && u.host === 'www.notion.so') {
                let hashBlockID = toDashID(item.hash.slice(1))
                item.href = `#${hashBlockID}`

                let block = document.querySelector(`div[data-block-id="${hashBlockID}"]`)
                if (block) {
                    block.id = hashBlockID
                }
            }
        });

        // bookmark 修复，notion更改了 bookmark block 的生成规则，a 标签内没有 href了
        document.querySelectorAll("#notion-app > div > div.notion-cursor-listener > div > div.notion-scroller.vertical.horizontal > div.notion-page-content > div[data-block-id] > div > div > a").forEach(a => {
            if (!a.href) {
                a.href = a.querySelector("div > div:first-child > div:last-child").innerText
            }
        })
        // 表格视图 CSS 修复
        document.querySelectorAll("div.notion-scroller.horizontal").forEach(item => {
            item.children[0].style.padding = 0
            item.previousElementSibling.style.paddingLeft = 0
            // 表格在 safari & edge 中显示有问题。
            item.style.overflowX = "scroll"
        })
        // 文章内容
        let content = document.querySelector('#notion-app > div > div.notion-cursor-listener > div > div > div.notion-page-content')


        // 可编辑内容修复
        let contenteditable = content.querySelectorAll("div[contenteditable=true]")
        contenteditable.forEach(i => {
            i.setAttribute("contenteditable", "false");
        })
        if (content) {
            return content.innerHTML
        }
        else {
            return 'error'
        }
    })

    await browser.close();
    return data;
}

async function doit() {
    const url = process.argv[2];
    const html = await getPageHtml(url);
    console.log(html);
}

doit();
