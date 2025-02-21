import http from 'k6/http'
import {Trend} from 'k6/metrics';
import grpc from 'k6/x/grpc';
import nacos from 'k6/x/nacos';
import builder from 'k6/x/grpc_builder';
import {check, sleep, group} from 'k6'
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

export function setup() {
    nacos.init("nacosClient1", "nacos.test.infra.ww5sawfyut0k.bitsvc.io", 8848, "bybit-nacos", "bybit-nacos", "efficiency-test");
    builder.fillBuilder("nacosClient1", "DEFAULT_GROUP")
}

const grpc_client_scenarioCode7425875b771a48e2b539079887ed199d_0 = new grpc.Client();
grpc_client_scenarioCode7425875b771a48e2b539079887ed199d_0.load('CpYCCgtoZWxsby5wcm90bxIFaGVsbG8iIgoMSGVsbG9SZXF1ZXN0EhIKBE5hbWUYASABKAlSBE5hbWUiTAoNSGVsbG9SZXNwb25zZRISCgRDb2RlGAEgASgJUgRDb2RlEicKBERhdGEYAiABKAsyEy5oZWxsby5SZXNwb25zZURhdGFSBERhdGEiNAoMUmVzcG9uc2VEYXRhEhIKBE5hbWUYASABKAlSBE5hbWUSEAoDQUdFGAIgASgJUgNBR0UyQAoFSGVsbG8SNwoIU2F5SGVsbG8SEy5oZWxsby5IZWxsb1JlcXVlc3QaFC5oZWxsby5IZWxsb1Jlc3BvbnNlIgBCDloMY2hhcmxlcy5zb25nYgZwcm90bzM=');
export const options = {
    discardResponseBodies: false,
    scenarios: {
        scenarioCode7425875b771a48e2b539079887ed199d: {
            duration: '600s',
            preAllocatedVUs: 1,
            rate: 1,
            executor: 'constant-arrival-rate',
            gracefulStop: '2s',
            maxVUs: 20,
            exec: 'scenarioExec0',
            timeUnit: '1s'
        }
    }
}

export function scenarioExec0() {
    group('hello.Hello/SayHello', function () {
        let SYS_TRACE_ID = gen_TraceId_SYS();
        if (__ITER == 0) {

        }
        // 'api-manage-service-efficiency-prod.test.efficiency.ww5sawfyut0k.bitsvc.io(http:80 grpc:9090)'
        // 'nacos:///api-manage-service'
        grpc_client_scenarioCode7425875b771a48e2b539079887ed199d_0.connect('api-manage-service-efficiency-prod.test.efficiency.ww5sawfyut0k.bitsvc.io:9090', {
            plaintext: true,
            shareConn: true
        });

        const res0 = grpc_client_scenarioCode7425875b771a48e2b539079887ed199d_0.invoke('hello.Hello/SayHello', {"Name": "aaa"}, {
            metadata: {
                'ss': 'ss',
                'traceparent': '00-' + SYS_TRACE_ID + '-' + gen_SpanId_SYS() + '-01'
            }
        });
        console.log(res0)
    });
}
