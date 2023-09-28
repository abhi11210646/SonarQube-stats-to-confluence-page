require("dotenv").config();
const sonar = require("./sonarstats");
const confluence = require("./confluence");
const { sonarConfig, confluenceConfig } = require("./config");


async function start() {
    try {
        // Generate stats
        const stats = [];
        for (let projectKey of sonarConfig.projectKeys) {
            console.log("Fetching stats from sonarQube for", projectKey);
            const data = await sonar.sonarStats(projectKey);
            stats.push(data);
        }
        console.log("Updating Confluence page with sonarQube stats");
        await confluence.updateByPageId(confluenceConfig.pageId, stats);
        console.log("Successfully updated Confluence page");

    } catch (error) {
        console.error("Error in fetching and updating stats");
        if (error.response) {
            console.log("Error data:", error.response.data);
            console.log("Status:", error.response.status);
        } else if (error.request) {
            console.log(error.request);
        } else {
            console.log('Error', error.message);
        }
    }
}


start();