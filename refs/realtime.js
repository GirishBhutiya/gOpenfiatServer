const fetch = require("node-fetch");
const { randomIntId } = require("../../random");
const { Fetcher } = require("../commons/commons");

module.exports.MakeLead = async function (req, ua) {
  let lead = { id: randomIntId(), ua, created: new Date().toISOString() };
  if (req.headers["x-real-ip"]) {
    lead.ip = req.headers["x-real-ip"];
  } else if (req.headers["x-forwarded-for"]) {
    lead.ip = req.headers["x-forwarded-for"];
  } else if (req.headers["x-vercel-forwarded-for"]) {
    lead.ip = req.headers["x-vercel-forwarded-for"];
  }
  lead.ip = "117.230.169.110";
  try {
    const promise = await fetch("https://ifconfig.co/json?ip=" + lead.ip);
    let data = await promise.json();
    if (data.ip) {
      lead = {
        ...lead,
        country: data.country || "",
        ccode: data.country_iso || "",
        region: data.region_name || "",
        rcode: data.region_code || "",
        city: data.city || "",
        loc: data.latitude + "," + data.longitude || "",
        asn: data.asn + " " + data.asn_org || "",
      };
    } else {
      const promise2 = await fetch(
        "https://ipinfo.io/" + lead.ip + "?token=119fd691d6a4b2"
      );
      data = await promise2.json();
      if (data.ip) {
        lead = {
          ...lead,
          country: "",
          ccode: data.country || "",
          region: data.region_name || "",
          rcode: "",
          city: data.city || "",
          loc: data.loc || "",
          asn: data.org || "",
        };
      }
    }
  } catch (error) {
    console.log(error);
  }

  //Get vercel headers
  /* if (req.headers["x-vercel-id"]) {
        data.misc = { reqid: req.headers["x-vercel-id"] };
      }
      if (req.headers) {
          data.misc.reqcountry = req.headers["x-vercel-ip-country"];
          data.misc.reqregion = req.headers["x-vercel-ip-country-region"];
          data.misc.reqcity = req.headers["x-vercel-ip-city"];
        } */

  return lead;
};
