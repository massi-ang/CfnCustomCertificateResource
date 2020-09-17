# Mirai::Iot::Certificate

A resource to create Certificate and Keys for IoT Things

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "Type" : "Mirai::Iot::Certificate",
    "Properties" : {
        "<a href="#bucket" title="Bucket">Bucket</a>" : <i>String</i>,
        "<a href="#status" title="Status">Status</a>" : <i>String</i>
    }
}
</pre>

### YAML

<pre>
Type: Mirai::Iot::Certificate
Properties:
    <a href="#bucket" title="Bucket">Bucket</a>: <i>String</i>
    <a href="#status" title="Status">Status</a>: <i>String</i>
</pre>

## Properties

#### Bucket

The S3 bucket name where to write the certificate

_Required_: Yes

_Type_: String

_Maximum_: <code>64</code>

_Pattern_: <code>^[a-zA-Z0-9_\-]{0,64}$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Status

The status of the certificate

_Required_: No

_Type_: String

_Allowed Values_: <code>ACTIVE</code> | <code>INACTIVE</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

## Return Values

### Ref

When you pass the logical ID of this resource to the intrinsic `Ref` function, Ref returns the Id.

### Fn::GetAtt

The `Fn::GetAtt` intrinsic function returns a value for a specified attribute of this type. The following are the available attributes and sample return values.

For more information about using the `Fn::GetAtt` intrinsic function, see [Fn::GetAtt](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html).

#### Arn

The Arn of this certificate

#### Id

The Id of this certificate

