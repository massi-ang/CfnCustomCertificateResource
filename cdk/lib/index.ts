import * as cdk from '@aws-cdk/core';
import * as s3 from '@aws-cdk/aws-s3';
import { countReset } from 'console';

export enum CertificateStatus {
  ACTIVE = "ACTIVE",
  INACITVE = "INACTIVE"
}
export interface IotCertificateProps {
  /**
   * The bucket where to store the key and certificate.
   */
  bucket: s3.Bucket,
  /**
   * The status of the certificate
   * 
   * @default(CertificateStatus.INACTIVE)
   */
  status?: CertificateStatus
}

export class IotCertificate extends cdk.Construct {
  /** @returns the ARN of the certificate
   * 
   */
  public readonly certificateArn: string;
  /** @returnsthe Id of the certificate
   * 
   */
  public readonly certificateId: string;

  constructor(scope: cdk.Construct, id: string, props: IotCertificateProps) {
    super(scope, id);

    const cert = new cdk.CfnResource(this, id+'_cert', {
      type: 'AWSSamples::Iot::Certificate',
      properties: {
        Bucket: props.bucket.bucketName,
        Status: props.status
      }
    })

    this.certificateArn = cert.getAtt('Arn').toString()
    this.certificateId = cert.getAtt('Id').toString()
  }
}
