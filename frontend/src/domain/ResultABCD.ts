import ow from 'ow'

export interface ResultABCD {
  name: string
  answers: number[]
}

export function validateResultABCD(
  o: ResultABCD | any,
): asserts o is ResultABCD {
  ow(o, 'ResultABCD', ow.object)
  ow(o.name, 'ResultABCD.name', ow.string.maxLength(40))
  ow(o.answers, 'ResultABCD.answers', ow.array.ofType(ow.number.finite.integer))
}
