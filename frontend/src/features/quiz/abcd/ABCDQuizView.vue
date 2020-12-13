<template>
  <b-row id="abcd-quiz-view">
    <div class="col-12">
      <state-matches :state="state">
        <template #Intro>
          <ABCDQuizIntro :interpreter="interpreter" :state="state" />
        </template>
        <template #Question>
          <ABCDQuizQuestion :interpreter="interpreter" :state="state" />
        </template>
        <template #ProvideDetails>
          <ABCDQuizProvideDetails :interpreter="interpreter" :state="state" />
        </template>
        <template #SavingResults>
          <Loading>Ładowanie wyników...</Loading>
        </template>
        <template #RedirectToSuccessPage>
          <ABCDQuizRedirectToResults
            :interpreter="interpreter"
            :state="state"
          />
        </template>
        <template #Error.NoResults>
          <Error>
            Wystąpił nieznany błąd. Jeśli chcesz nam pomóc rozwijać Siejmy, to
            powiadom nas o tym!
          </Error>
        </template>
        <template #Error.HasResults>
          <ABCDQuizTemporaryResults :interpreter="interpreter" :state="state" />
        </template>
      </state-matches>
    </div>
  </b-row>
</template>

<script lang="ts">
// tslint:disable member-ordering
import { Error, Loading, StateMatches } from '@/components'
import { QuizABCD } from '@/domain'
import { ResultsRepositoryAjax } from '@/services'
import { Component, Inject, Prop, Vue } from 'vue-property-decorator'

import {
  ABCDQuizIntro,
  ABCDQuizProvideDetails,
  ABCDQuizQuestion,
  ABCDQuizRedirectToResults,
  ABCDQuizTemporaryResults,
} from './components'
import { interpretMachine } from './interpreter'
import { ABCDQuizInterpreter } from './machine'

@Component({
  components: {
    StateMatches,
    Loading,
    Error,
    ABCDQuizIntro,
    ABCDQuizQuestion,
    ABCDQuizProvideDetails,
    ABCDQuizRedirectToResults,
    ABCDQuizTemporaryResults,
  },
})
export default class extends Vue {
  @Prop({ required: true, type: String })
  public quizUrl!: string

  @Prop({ required: true, type: Object })
  public quiz!: QuizABCD

  @Inject()
  private resultRepository!: ResultsRepositoryAjax

  public interpreter: ABCDQuizInterpreter = interpretMachine({
    resultRepository: this.resultRepository!,
    quiz: this.quiz,
  })
  public state: ABCDQuizInterpreter['state'] = this.interpreter.initialState

  public created() {
    this.startMachine()
  }

  public beforeDestroy() {
    this.interpreter.stop()
  }

  private startMachine() {
    this.interpreter
      .onTransition(state => {
        this.state = state
      })
      .start()
  }
}
</script>
