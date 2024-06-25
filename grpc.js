import grpc from 'k6/x/grpc';
import nacos from "k6/x/nacos";
import builder from 'k6/x/grpc_builder';

export function setup() {
    nacos.init("testNacos", "nacos.test.infra.ww5sawfyut0k.bitsvc.io", 8848, "nacos", "nacos", "byone-auto-test");
    nacos.init("publicNacos", "nacos.test.infra.ww5sawfyut0k.bitsvc.io", 8848, "nacos", "nacos", "public");
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
            duration: '300s',
            // How many iterations per timeUnit
            rate: 10,
            // Start `rate` iterations per second
            timeUnit: '1s',
            // Pre-allocate VUs
            preAllocatedVUs: 20,
            exec: 'scenarioExec0',
        }
    }
}

export function scenarioExec0() {
    const addr = "10.110.65.20:9008";
    //开启反射：reflect:true
    // grpc_client_scenarioCode92fb22b8775544cca7104c88c08e7c8e_0.connect("nacos:///award-srv", {plaintext: true,timeout:'3s'});
    // const response = grpc_client_scenarioCode92fb22b8775544cca7104c88c08e7c8e_0.invoke('grpc.health.v1.Health/Watch', {"service": "AwardingSrv"});
    // console.log(response);


    if (__ITER == 0) {
        // grpc_client_scenarioCode92fb22b8775544cca7104c88c08e7c8e_0.connect("10.110.113.150:9205", {
        //     plaintext: true,
        //     timeout: '3s'
        // });

        grpc_client_scenarioCode92fb22b8775544cca7104c88c08e7c8e_0.connect("nacos:///byone.greeter.rpc.local", {
            plaintext: true,
            timeout: '10s'
        });
        console.log(__VU);
    }

    const response = grpc_client_scenarioCode92fb22b8775544cca7104c88c08e7c8e_0.invoke('com.bybit.infra.test.bsoatest.Greeter/SayHello', {
        "name": "haha"
    }, {
        metadata: {
            'user-id': "1"
        },
        timeout: '5s'
    });


    console.log(__VU,response);
}