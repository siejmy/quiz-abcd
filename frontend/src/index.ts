import { BootstrapVue } from 'bootstrap-vue'
import Vue from 'vue'
import VueRouter from 'vue-router'

import App from './App.vue'
import { Configuration, validateConfiguration } from './Configuration'
import { QuizPage } from './features/quiz'
import { ResultsRepositoryAjax } from './services'

Vue.config.productionTip = false
Vue.use(VueRouter)
Vue.use(BootstrapVue)

export function mountQuiz(tag: string, config: Configuration) {
  validateConfiguration(config)
  const routes = [
    { path: '/', component: QuizPage },
    { path: '*', component: QuizPage },
  ]
  new Vue({
    render: h => h(App),
    router: new VueRouter({
      mode: 'hash',
      routes,
    }),
    provide: {
      resultRepository: new ResultsRepositoryAjax(config.saveUrl),
    },
    data: {
      config,
    },
  }).$mount(tag)
}
