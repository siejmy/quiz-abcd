<template>
  <div id="quiz-redirect-to-results">
    <p class="redirect-to-results">
      <strong>Przekierowywanie do wyniku...</strong>
      Jeśli wynik nie otwiera się automatycznie
      <a href="resultUrl">kliknij tu, aby zobaczyć wynik</a>.
    </p>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'

import { ABCDQuizInterpreter } from '../machine'

@Component({
  components: {},
})
export default class extends Vue {
  @Prop({ required: true })
  public interpreter!: ABCDQuizInterpreter

  @Prop({ required: true })
  public state!: ABCDQuizInterpreter['state']

  get resultUrl(): string {
    return this.state.context.resultUrl
  }

  public mounted() {
    window.location.href = this.resultUrl
  }
}
</script>
<style scoped>
.redirect-to-results {
  color: white;
}
</style>
