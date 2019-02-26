const PROXY_CONFIG = [
  {
    context: [
      '/api/v1/checks',
    ],
    target: 'http://52.51.26.115:32641/api/v1/checks',
    secure: false
  },
  {
    context: [
      '/api/v1/plugins',
    ],
    target: 'http://52.51.26.115:32641/api/v1/plugins',
    secure: false
  },
  {
    context: [
      '/check'
    ],
    target: 'http://52.51.26.115:32641/check',
    secure: false
  }
];

module.exports = PROXY_CONFIG;
