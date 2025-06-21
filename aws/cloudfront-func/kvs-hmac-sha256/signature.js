import crypto from 'crypto';

const secret = 'godev';

const signingMethod = 'sha256';
const response403InvalidSign = {
    statusCode: 403,
    statusDescription: 'Unauthorized invalid get header signature'
};

const response403InvalidExpiredTime = {
    statusCode: 403,
    statusDescription: 'Unauthorized invalid get header expired'
};

const response403NotMatch = {
    statusCode: 403,
    statusDescription: 'Unauthorized header signature not match'
};

const response403Expired = {
    statusCode: 403,
    statusDescription: 'Request expired'
};


async function handler(event) {
    var requestObj = event.request;
    var queryString = requestObj.querystring;

    // console.log(requestObj)

    var signature = queryString.signature.value;
    // console.log(signature)
    if (!signature) {
        return response403InvalidSign
    }

    var expiredTime = queryString.expiredTime.value
    // console.log(expiredTime)
    if (!expiredTime) {
        return response403InvalidExpiredTime
    }

    var sign = crypto.createHmac(signingMethod, secret).update(expiredTime).digest('hex');
    if (sign !== signature) {
        return response403NotMatch
    }

    const now = Math.floor(Date.now() / 1000);
    console.log(now)
    console.log(expiredTime)
    if (now > expiredTime){
        return response403Expired
    }
    return requestObj
}
