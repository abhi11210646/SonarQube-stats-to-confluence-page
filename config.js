
const sonarConfig = {
    "host": "https://sonarqube.one.com",
    "api_key": process.env.SONAR_API_KEY,
    "projectKeys": [
        "app.webmail",
        "CompanionApp",
        "Webshop",
        "one.com-wp-addons-assets"
    ],
    "metricKeys": [
        "code_smells",
        "critical_severity_vulns",
        "bugs",
        "alert_status",
        "quality_gate_details"
    ]
}
const confluenceConfig = {
    "host": "https://group-one.atlassian.net/wiki/rest/",
    "api_key": process.env.CONFLUENCE_API_KEY,
    "pageId": 32589873205
}


module.exports = { confluenceConfig, sonarConfig }