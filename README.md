# twitchstream is a minimalistic twitch app

## Shows stream, chat and last events

Try it [here](https://stream.seregayoga.com)

Run locally:
`docker-compose up`

It can be deployed and scaled easily on almost every cloud using load balancer (ELB) and identical instances with the app on VMs (EC2), because it's completely stateless (including session management).

Authorization code snipppets were borrowed from https://github.com/twitchdev/authentication-samples
