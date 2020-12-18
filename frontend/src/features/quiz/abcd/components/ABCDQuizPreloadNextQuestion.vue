<template>
  <PreloadImage v-if="imgToPreloadSrc" :src="imgToPreloadSrc" />
</template>

<script lang="ts">
import { PreloadImage } from '@/components'
import { QuizABCDQuestion } from '@/domain'
import { Component, Prop, Vue } from 'vue-property-decorator'

import { ABCDQuizInterpreter } from '../machine'

@Component({
  components: {
    PreloadImage,
  },
})
export default class extends Vue {
  @Prop({ required: true })
  public state!: ABCDQuizInterpreter['state']

  public get isIntro(): boolean {
    return this.state.matches('Intro')
  }

  public get currentQuestionNo(): number {
    return this.state.context.currentQuestionNo
  }

  public get totalQuestions() {
    return this.state.context.quiz.questions.length
  }

  public get questionNoToPreload(): number | undefined {
    if (this.isIntro) {
      return 0
    }

    if (this.currentQuestionNo + 1 < this.totalQuestions) {
      return this.currentQuestionNo + 1
    }
  }

  public get questionToPreload(): QuizABCDQuestion | undefined {
    if (this.questionNoToPreload === undefined) {
      return
    }
    return this.state.context.quiz.questions[this.questionNoToPreload]
  }

  public get imgToPreloadSrc(): string | undefined {
    if (this.questionToPreload) {
      return this.questionToPreload.imageUrl
    }
  }
}
</script>
