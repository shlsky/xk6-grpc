import grpc from 'k6/x/grpc';
import nacos from "k6/x/nacos";
import builder from 'k6/x/grpc_builder';
import { check, sleep, group } from 'k6'
import jsonpath from 'https://jslib.k6.io/jsonpath/1.0.2/index.js';



export function setup() {
    nacos.init("testNacos", "nacos.io", 8848, "nacos", "nacos", "test");
    nacos.init("publicNacos", "nacos.io", 8848, "nacos", "nacos", "public");
    builder.fillBuilder("testNacos", "")
}


const grpc_client_scenarioCode92fb22b8775544cca7104c88c08e7c8e_0 = new grpc.Client();
grpc_client_scenarioCode92fb22b8775544cca7104c88c08e7c8e_0.load('CsADChBoZWxsb3dvcmxkLnByb3RvEh1jb20uYnliaXQuaW5mcmEudGVzdC5ic29hdGVzdCIiCgxIZWxsb1JlcXVlc3QSEgoEbmFtZRgBIAEoCVIEbmFtZSImCgpIZWxsb1JlcGx5EhgKB21lc3NhZ2UYASABKAlSB21lc3NhZ2Uy2gEKB0dyZWV0ZXISZAoIU2F5SGVsbG8SKy5jb20uYnliaXQuaW5mcmEudGVzdC5ic29hdGVzdC5IZWxsb1JlcXVlc3QaKS5jb20uYnliaXQuaW5mcmEudGVzdC5ic29hdGVzdC5IZWxsb1JlcGx5IgASaQoNU2F5SGVsbG9BZ2FpbhIrLmNvbS5ieWJpdC5pbmZyYS50ZXN0LmJzb2F0ZXN0LkhlbGxvUmVxdWVzdBopLmNvbS5ieWJpdC5pbmZyYS50ZXN0LmJzb2F0ZXN0LkhlbGxvUmVwbHkiAEJcCiNjb20uYnliaXQuaW5mcmEudGVzdC5ic29hdGVzdC5wcm90b0IOR3JlZXRlclNlcnZpY2VQAVojY29tLmJ5Yml0LmluZnJhLnRlc3QuYnNvYXRlc3QucHJvdG9iBnByb3RvMw==');
export const options = {
    discardResponseBodies: false,
    scenarios: {
        scenarioCode92fb22b8775544cca7104c88c08e7c8e: {
            executor: 'constant-arrival-rate',
            // How long the test lasts
            duration: '120s',
            // How many iterations per timeUnit
            rate: 100,
            // Start `rate` iterations per second
            timeUnit: '1s',
            // Pre-allocate VUs
            preAllocatedVUs: 100,
            exec: 'scenarioExec0',
        }
    }
}

export function scenarioExec0() {

    if (__ITER == 0) {
        // grpc_client_scenarioCode92fb22b8775544cca7104c88c08e7c8e_0.connect("10.110.113.150:9205", {
        //     plaintext: true,
        //     timeout: '3s'
        // });

        grpc_client_scenarioCode92fb22b8775544cca7104c88c08e7c8e_0.connect("nacos:///rpc.local", {
            plaintext: true,
            shareConn: true,
            timeout: '10s'
        });
        console.log(__VU);
    }

    const res0 = grpc_client_scenarioCode92fb22b8775544cca7104c88c08e7c8e_0.invoke('com.Greeter/SayHello', {
        "name": "haha"
    }, {
        metadata: {
            'user-id': "1"
        },
        timeout: '5s'
    });
    console.log(res0.message)
}