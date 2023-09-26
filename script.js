const sonar = require("./sonarstats");
const confluence = require("./confluence");
const { sonarConfig, confluenceConfig } = require("./credentials.json");


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

    } catch (err) {
        console.error("Error in fetch and updating stats", err);
    }

}


start();