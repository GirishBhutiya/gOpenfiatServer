import { parseCookies, setCookie } from "nookies";

export default async function Realtime() {
  try {
    const cookies = parseCookies();
    if (cookies["LERNEN_AD"] && cookies["LERNEN_AD"].length > 0) {
      return cookies["LERNEN_AD"];
    } else {
      let ua = {};
      if (navigator) {
        if (navigator.platform) {
          let p = navigator.platform;
          ua.os = navigator.platform;
          if (p.indexOf("Win") != -1) ua.os = "Windows";
          if (p.indexOf("Mac") != -1) ua.os = "MacOS";
          if (p.indexOf("iPhone") != -1) ua.os = "iOS";
          if (p.indexOf("Android") != -1) ua.os = "Android";
          if (p.indexOf("Linux") != -1) ua.os = "Linux";
          if (p.indexOf("X11") != -1) ua.os = "UNIX";
        }
        if (navigator.userAgent) {
          const u = navigator.userAgent;
          ua.type = u;
          if (u.includes("Mobile")) {
            ua.type = "Mobile";
          } else {
            ua.type = "Desktop";
          }
          ua.browser = u;
          if (
            u.includes("Safari") &&
            (u.includes("Mac") || u.includes("iPhone"))
          ) {
            //version u.split(" ").find(i=>i.includes("Safari/")).replace("Safari/","")
            ua.browser = "Apple Safari";
          } else if (u.includes("Opera") || u.includes("OPR")) {
            ua.browser = "Opera";
          } else if (u.includes("Edge") || u.includes("Edg")) {
            ua.browser = "Microsoft Edge";
          } else if (u.includes("MSIE")) {
            ua.browser = "Internet Explorer";
          } else if (u.includes("Firefox") || u.includes("FxiOS")) {
            ua.browser = "Mozilla Firefox";
          } else if (navigator.brave) {
            ua.browser = "Brave Browser";
          } else if (u.includes("Chromium")) {
            ua.browser = "Chromium";
          } else if (u.includes("Chrome") || u.includes("CriOS")) {
            ua.browser = "Google Chrome";
          } else {
            ua.browser = "Unknown Browser";
          }
        }
      }
      try {
        const promise = await fetch(
          "/api/realtime?ua=" + encodeURIComponent(JSON.stringify(ua))
        );
        const { success, ld } = await promise.json();
        if (success) {
          console.log("Realtime request complete.");
          setCookie(null, "LERNEN_AD", ld);
          return ld;
        } else {
          console.log("Realtime update failed.");
        }
      } catch (error) {
        console.log(error);
        console.log("Realtime request failed.");
      }
    }
  } catch (error) {
    console.log(error);
    console.log("Realtime request failed.");
  }
  return null;
}
