<p>

Packages:

</p>

<ul>

<li>

<a href="#argoproj.io%2fv1alpha1">argoproj.io/v1alpha1</a>

</li>

</ul>

<h2 id="argoproj.io/v1alpha1">

argoproj.io/v1alpha1

</h2>

<p>

<p>

Package v1alpha1 is the v1alpha1 version of the API.

</p>

</p>

Resource Types:

<ul>

</ul>

<h3 id="argoproj.io/v1alpha1.HTTPSubscription">

HTTPSubscription

</h3>

<p>

(<em>Appears on:</em>
<a href="#argoproj.io/v1alpha1.SubscriptionSpec">SubscriptionSpec</a>)

</p>

<p>

<p>

HTTPSubscription describes the subscription details over HTTP

</p>

</p>

<table>

<thead>

<tr>

<th>

Field

</th>

<th>

Description

</th>

</tr>

</thead>

<tbody>

<tr>

<td>

<code>name</code></br> <em> string </em>

</td>

<td>

</td>

</tr>

<tr>

<td>

<code>url</code></br> <em> string </em>

</td>

<td>

</td>

</tr>

</tbody>

</table>

<h3 id="argoproj.io/v1alpha1.NATSSubscription">

NATSSubscription

</h3>

<p>

(<em>Appears on:</em>
<a href="#argoproj.io/v1alpha1.SubscriptionSpec">SubscriptionSpec</a>)

</p>

<p>

<p>

NATSSubscription describes the subscription details over NATS protocol

</p>

</p>

<table>

<thead>

<tr>

<th>

Field

</th>

<th>

Description

</th>

</tr>

</thead>

<tbody>

<tr>

<td>

<code>name</code></br> <em> string </em>

</td>

<td>

<p>

Name of the subscription

</p>

</td>

</tr>

<tr>

<td>

<code>serverURL</code></br> <em> string </em>

</td>

<td>

<p>

ServerURL is NATS server URL

</p>

</td>

</tr>

<tr>

<td>

<code>subject</code></br> <em> string </em>

</td>

<td>

<p>

Subject is the name of the NATS subject

</p>

</td>

</tr>

</tbody>

</table>

<h3 id="argoproj.io/v1alpha1.Subscription">

Subscription

</h3>

<p>

<p>

Subscription is the definition of a subscription resource

</p>

</p>

<table>

<thead>

<tr>

<th>

Field

</th>

<th>

Description

</th>

</tr>

</thead>

<tbody>

<tr>

<td>

<code>metadata</code></br> <em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta </a> </em>

</td>

<td>

Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.

</td>

</tr>

</tbody>

</table>

<h3 id="argoproj.io/v1alpha1.SubscriptionSpec">

SubscriptionSpec

</h3>

<p>

<p>

SubscriptionSpec describes the specification of the subscription
resource

</p>

</p>

<table>

<thead>

<tr>

<th>

Field

</th>

<th>

Description

</th>

</tr>

</thead>

<tbody>

<tr>

<td>

<code>http</code></br> <em>
<a href="#argoproj.io/v1alpha1.HTTPSubscription"> \[\]HTTPSubscription
</a> </em>

</td>

<td>

<p>

HTTP refers to list of subscriptions over HTTP protocol

</p>

</td>

</tr>

<tr>

<td>

<code>nats</code></br> <em>
<a href="#argoproj.io/v1alpha1.NATSSubscription"> \[\]NATSSubscription
</a> </em>

</td>

<td>

<p>

NATS refers to list of subscriptions over NATS protocol

</p>

</td>

</tr>

</tbody>

</table>

<h3 id="argoproj.io/v1alpha1.SubscriptionStatus">

SubscriptionStatus

</h3>

<p>

<p>

SubscriptionStatus describes the status of the subscription resource

</p>

</p>

<table>

<thead>

<tr>

<th>

Field

</th>

<th>

Description

</th>

</tr>

</thead>

<tbody>

<tr>

<td>

<code>createdAt</code></br> <em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#time-v1-meta">
Kubernetes meta/v1.Time </a> </em>

</td>

<td>

<p>

CreatedAt refers to creation time

</p>

</td>

</tr>

<tr>

<td>

<code>updatedAt</code></br> <em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#time-v1-meta">
Kubernetes meta/v1.Time </a> </em>

</td>

<td>

<p>

UpdatedAt refers to time at the resource was updated

</p>

</td>

</tr>

</tbody>

</table>

<hr/>

<p>

<em> Generated with <code>gen-crd-api-reference-docs</code> on git
commit <code>2bdb384</code>. </em>

</p>
