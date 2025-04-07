import crypto from 'crypto';

const secret = 'godev';

const signingMethod = 'sha256';
const response401InvalidSign = {
    statusCode: 401,
    statusDescription: 'Unauthorized invalid get header signature'
};

const response401InvalidRequestId = {
    statusCode: 401,
    statusDescription: 'Unauthorized invalid get header requestid'
};


const response401NotMatch = {
    statusCode: 401,
    statusDescription: 'Unauthorized header signature not match'
};

const response200 = {
    statusCode: 200,
    statusDescription: 'success'
};

async function handler(event) {
    var requestObj = event.request;
    var headers = requestObj.headers;

    var signature = headers.signature.value;
    if (!signature) {
        return response401InvalidSign
    }

    var requestId = headers.requestid.value
    if (!requestId) {
        return response401InvalidRequestId
    }

    var sign = crypto.createHmac(signingMethod, secret).update(requestId).digest('hex');
    if (sign !== signature) {
        return response401NotMatch
    }
    return response200
}