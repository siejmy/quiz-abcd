import { ResultABCD, validateResultABCD } from '@/domain'
import ow from 'ow'
export class ResultsRepositoryAjax {
  constructor(private saveUrl: string) {
    return
  }

  public async saveResult(r: ResultABCD): Promise<{ id: string; url: string }> {
    validateResultABCD(r)
    const response = await fetch(this.saveUrl, {
      method: 'POST',
      cache: 'no-cache',
      body: JSON.stringify(r),
    })
    const responseData = await response.json()
    ow(responseData, 'responseData', ow.object)
    ow(responseData.id, 'responseData.id', ow.string.nonEmpty)
    ow(responseData.url, 'responseData.url', ow.string.nonEmpty.url)
    return responseData
  }
}
