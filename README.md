election
========

A start at implementing an online minimum trust secret ballot election system.


Background
==========

An `Election`, as considered herein, consists of collecting `ballot`s from a subset of `voters`, and publishing the results.

Basically:
- all `voters` can (but are not required) to `vote` meaning contribute exactly one `ballot` to the `election`.
- the set of `ballot`s contributed to the election publicly visible upon the end of the `election`

There are some additional goals:
- resistance to coercion: it is impossible to force a voter to vote a particular way
- verifiability:
    - (1) voters can guarantee that their `ballot` is included and unedited
    - (2) voters can guarantee that each voter contributes at most one `ballot`
    - (3) voters can guarantee no `ballot`s are included without being contributed by a voter
    - (4) voters can guarantee the set of voters is correct (no one denied, no one wrongly allowed to vote)
    
Unfortunately, in any realistic situation this is impossible. Perfect resistance to coercion requires
voters to not be able to be controlled in any way,
and verifiability requires at minimum trust that the list of potential voters ia correct.

Despite these issues, elections are conducted, and in some cases, are done well enough that the public accepts the results.
My goal for this project is to do at better than existing systems
(particularly paper ballot, and existing electronic systems in the United States), while also lowering the cost.


Particularly, I'll target beating the vote by mail systems, since they seem to be well accepted,
are easy to analyze, and really easy to beat.
This means that this system aims to better resistance to coercion and higher verifiability
than existing vote by mail systems.


Coercion issues
+++++++++++++++

The first defense against coercion is the secret ballot: the property that one can not associate a ballot with its voter.
Without this, its easy to reward or punish people for voting a particular way. This can fail for trivially small numbers of voters,
but generally seems to be accepted as effective by the public.

The second defense generally seems on the legal side: in some countries, you will get in very serious trouble if you are found
trying to coerce votes.

Its important to note that one should not be able to share their ballot in a manner that can prove that they really voted or will vote that way,
since then people could simply be coerced to do that. This is what voting booths are all about, and this is the main place vote by mail fails.
I'll mention here that my system has about the same issues as vote my mail here (perhaps worse),
but could also be done in a voting booth type setup which would fix the issue. More about that later.

This is a particularly hard problem to solve in a robust manner since it includes vulnerabilities to social engineering.
Threats, bribes, replacing people with replicants, mind control etc. are all attack vectors here. Its not a solvable problem,
but hopefully we can prevent it from being a significant issue in actual elections. When analyzing an electoral system, always
keep coercion in mind, since it will potentially be an issue, you need to evaluate how bad of an issue it is!

Verifiability issues
++++++++++++++++++++

This is where the existing system fail really badly. Of the three factors listed above, existing systems generally
don't even to a partial job of any of them. Its purely done based on trusting the administration of the election.
This is very bad, especially given the elections are often not run by trusted neutral parties.

Of the 4 factors listed above, my system solves (1), and does at least as good at the others.

There are also some nasty issues regarding enforcing when the polls should close. I'm pretty sure this isn't a perfect solution to that,
but its pretty easy to come up with something decent. I we have to ask voters to accept that corrupt management of the ballot server could
introduce some fraud if people are still trying to vote right up to and/or after the deadline regarding which of these votes are rejected.

The Design
==========

- A document (the `election description`) listing the public keys for the various involved servers, and any other needed meta dada (like what is being voted on)
is published
- A set of public keys for registered voters is published. This does not need to (and likely should not be) anonymized.
- Equivalent but anonymized sets of public keys are constructed and published on a series of `AnonymizationServer`s.
- A set of ballots signed with the anonymized public keys is constructed and published.

All communication is stateless, and all messages are idempotent. The origin of any requests not considered.
All wellformed requests are replied to with either proof that the request should not be met, or sufficient evidence 
to prove the server has committed fraud if the request is not met. 

Ex: If you submit a ballot, signed with your anonymized private+public key pair, and a signed message from the last AnonymizationServer
claiming your key is valid, a the ballot server must respond with one of:
- a copy of your message, signed with its own key. This means if the final ballot list does not contain your vote,
you can publish this response, which proves the ballot server signed, but did not include, a ballot, which is fraud.
- another different ballot, signed with your private key. This is proof that your request to add this ballot is invalid.

If you do not get one of these responses, you can publish the request you sent to the ballot server for others (such as auditors, election monitors, media, etc)
to try. Either they will get invalid responses (proving the fraud),
or they will get valid responses which means they can verify the server behaved correctly and has no performed the requested action.

