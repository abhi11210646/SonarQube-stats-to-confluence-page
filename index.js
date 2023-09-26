const axios = require('axios');

const axiosConfig = {
    headers: {
        'Authorization': `Basic ${Buffer.from(
            'XX:YY'
        ).toString('base64')}`,
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    }
};

const apiEndpoint = "https://group-one.atlassian.net/wiki/rest/api/content/32588300309?expand=body.storage";


const bodyData = {
    "version": {
        "number": "4", "message": "heyyyy"
    },
    "type": "page",
    "title":"hello tile222",
    "body": {
        "storage": {
            "value": '<p>this is body222222</p><p /><table data-table-width="760" data-layout="default" ac:local-id="447fa2a9-902a-4e64-8e90-e18acfcbc396"><tbody><tr><th><p><strong>ss</strong></p></th><th><p><strong>ww</strong></p></th><th><p><strong>qq</strong></p></th></tr><tr><td><p>3</p></td><td><p>4</p></td><td><p>7</p></td></tr><tr><td><p>2</p></td><td><p>22</p></td><td><p>11</p></td></tr></tbody></table><p>this is body</p><p /><p /><p>this is body</p><p /><p /><p>this is body</p><p /><p />',
            "representation": 'storage'
        }
    }
};


axios.put(apiEndpoint, bodyData, axiosConfig)
    .then((response) => {
        console.log(response.data);
    })
    .catch((error) => {
        console.error('=>>>>>> :', error.response.data);
    });