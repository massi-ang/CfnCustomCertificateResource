import { expect as expectCDK, haveResource } from '@aws-cdk/assert';
import * as cdk from '@aws-cdk/core';
import * as IotCertificate from '../lib/index';
import * as s3 from '@aws-cdk/aws-s3';

test('Certificate created', () => {
    const app = new cdk.App();
    const stack = new cdk.Stack(app, "TestStack");
    // WHEN
  const b = new s3.Bucket(stack, "test-bucket");
  new IotCertificate.IotCertificate(stack, 'MyTestConstruct', {
      bucket: b
    });
    // THEN
    expectCDK(stack).to(haveResource("AWSSample::Iot::Certificate"));
});

