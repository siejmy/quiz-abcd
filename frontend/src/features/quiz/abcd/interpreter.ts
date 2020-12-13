import { QuizABCD } from '@/domain'
import { ResultsRepositoryAjax } from '@/services'
import { interpret } from 'xstate'

import { ABCDQuizInterpreter, abcdQuizMachine, initialContext } from './machine'

export function interpretMachine({
  resultRepository,
  quiz,
}: {
  resultRepository: ResultsRepositoryAjax
  quiz: QuizABCD
}): ABCDQuizInterpreter {
  return interpret(
    abcdQuizMachine
      .withConfig({
        services: {
          saveResults: async ctx => {
            return resultRepository.saveResult(ctx.resultData)
          },
        },
      })
      .withContext({
        ...initialContext,
        quiz,
      }),
  )
}
