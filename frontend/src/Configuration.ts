import ow from 'ow'

import { QuizABCD, validateQuizABCD } from './domain'

export interface Configuration {
  title: string
  saveUrl: string
  reCaptchaKey: string
  facebookAppId: string
  quiz: QuizABCD
}

export function validateConfiguration(c: Configuration) {
  ow(c, 'Configuration', ow.object)
  ow(c.saveUrl, 'Configuration.saveUrl', ow.string.nonEmpty)
  ow(c.title, 'Configuration.title', ow.string.nonEmpty)
  ow(c.reCaptchaKey, 'Configuration.reCaptchaKey', ow.string.nonEmpty)
  ow(c.facebookAppId, 'Configuration.facebookAppId', ow.string.nonEmpty)

  validateQuizABCD(c.quiz)
}
