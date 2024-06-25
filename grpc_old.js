import grpc from 'k6/x/grpc';
import file from 'k6/x/file';
import {check, sleep, group} from 'k6'
import {SharedArray} from 'k6/data'
import papaparse from 'https://jslib.k6.io/papaparse/5.1.1/index.js';
import {URLSearchParams} from 'https://jslib.k6.io/url/1.0.0/index.js';
import jsonpath from 'https://jslib.k6.io/jsonpath/1.0.2/index.js';
import {randomIntBetween} from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import {uuidv4} from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

const Variables = {};
const util = new grpc.Util();

function gen_TraceId_SYS() {
    const randomUUID = uuidv4();
    return randomUUID.toString().replace(/-/g, "");
}

function gen_SpanId_SYS() {
    return gen_TraceId_SYS().slice(0, 16);
}

const grpc_client_scenarioCode38b8aa89aebc4123b9af94a85af25110_0 = new grpc.Client();
grpc_client_scenarioCode38b8aa89aebc4123b9af94a85af25110_0.load('CsADChBoZWxsb3dvcmxkLnByb3RvEh1jb20uYnliaXQuaW5mcmEudGVzdC5ic29hdGVzdCIiCgxIZWxsb1JlcXVlc3QSEgoEbmFtZRgBIAEoCVIEbmFtZSImCgpIZWxsb1JlcGx5EhgKB21lc3NhZ2UYASABKAlSB21lc3NhZ2Uy2gEKB0dyZWV0ZXISZAoIU2F5SGVsbG8SKy5jb20uYnliaXQuaW5mcmEudGVzdC5ic29hdGVzdC5IZWxsb1JlcXVlc3QaKS5jb20uYnliaXQuaW5mcmEudGVzdC5ic29hdGVzdC5IZWxsb1JlcGx5IgASaQoNU2F5SGVsbG9BZ2FpbhIrLmNvbS5ieWJpdC5pbmZyYS50ZXN0LmJzb2F0ZXN0LkhlbGxvUmVxdWVzdBopLmNvbS5ieWJpdC5pbmZyYS50ZXN0LmJzb2F0ZXN0LkhlbGxvUmVwbHkiAEJcCiNjb20uYnliaXQuaW5mcmEudGVzdC5ic29hdGVzdC5wcm90b0IOR3JlZXRlclNlcnZpY2VQAVojY29tLmJ5Yml0LmluZnJhLnRlc3QuYnNvYXRlc3QucHJvdG9iBnByb3RvMw==');
export const options = {
    discardResponseBodies: false,
    scenarios: {
        scenarioCode38b8aa89aebc4123b9af94a85af25110: {
            duration: '600s',
            preAllocatedVUs: 200,
            rate: 10,
            executor: 'constant-arrival-rate',
            gracefulStop: '2s',
            maxVUs: 200,
            exec: 'scenarioExec0',
            timeUnit: '1s'
        }
    }
}

export function scenarioExec0() {
    group('com.bybit.infra.test.bsoatest.Greeter/SayHello', function () {
        let SYS_TRACE_ID = gen_TraceId_SYS();
        if (__ITER == 0) {
            grpc_client_scenarioCode38b8aa89aebc4123b9af94a85af25110_0.connect('10.110.113.150:9205', {plaintext: true});
        }
        const res0 = grpc_client_scenarioCode38b8aa89aebc4123b9af94a85af25110_0.invoke('com.bybit.infra.test.bsoatest.Greeter/SayHello', {"name": "a"}, {metadata: {'traceparent': '00-' + SYS_TRACE_ID + '-' + gen_SpanId_SYS() + '-01'}});
        // file.appendString('/home/sys_su/pts/k6/6979/report/logs.log', SYS_TRACE_ID + ',' + 'com.bybit.infra.test.bsoatest.Greeter/SayHello' + ',' + res0.duration + ',' + JSON.stringify(res0.error) + '\n')
    });
}