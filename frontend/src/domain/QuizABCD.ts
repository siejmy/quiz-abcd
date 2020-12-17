import ow from 'ow'

export interface QuizABCD {
  id: string
  type: 'abcd'
  title: string
  introHtml: string
  introLegal?: string
  introImageUrl: string
  questions: QuizABCDQuestion[]
}

export interface QuizABCDQuestion {
  title: string
  imageUrl?: string
  distractors: string[]
  correctNo: number
  legal?: string
}

export function validateQuizABCD(o: QuizABCD | any): asserts o is QuizABCD {
  ow(o, 'QuizABCD', ow.object)
  ow(o.id, 'QuizABCD.id', ow.string.nonEmpty)
  ow(o.title, 'QuizABCD.title', ow.string.nonEmpty)
  ow(o.introHtml, 'QuizABCD.introHtml', ow.string.nonEmpty)
  ow(o.introImageUrl, 'QuizABCD.introImageUrl', ow.string.nonEmpty)
  ow(o.introLegal, 'QuizABCD.introLegal', ow.any(ow.string, ow.undefined))
  ow(o.type, 'QuizABCD.type', ow.string.equals('abcd'))
  ow(o.questions, 'QuizABCD.questions', ow.array.ofType(ow.object))
  o.questions.forEach((q: any) => validateQuizABCDQuestion(q))
}

export function validateQuizABCDQuestion(
  o: QuizABCDQuestion | any,
): asserts o is QuizABCDQuestion {
  ow(o, 'QuizABCDQuestion', ow.object)
  ow(o.title, 'QuizABCDQuestion.title', ow.string.nonEmpty)
  ow(
    o.imageUrl,
    'QuizABCDQuestion.imageUrl',
    ow.any(ow.undefined, ow.string.nonEmpty),
  )
  ow(
    o.distractors,
    'QuizABCDQuestion.distractors',
    ow.array.ofType(ow.string.nonEmpty),
  )
  ow(o.correctNo, 'QuizABCDQuestion.correctNo', ow.number.integer.finite)
  ow(o.legal, 'QuizABCDQuestion.introLegal', ow.any(ow.string, ow.undefined))
}
